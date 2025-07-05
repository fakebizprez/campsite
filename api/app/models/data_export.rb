# frozen_string_literal: true

class DataExport < ApplicationRecord
  include PublicIdGenerator
  include ImgixUrlBuilder
  include Rails.application.routes.url_helpers

  belongs_to :member, class_name: "OrganizationMembership"
  belongs_to :subject, polymorphic: true
  has_many :resources, class_name: "DataExportResource", dependent: :destroy

  def completed?
    completed_at.present?
  end

  def perform
    create_resources
    queue_resource_jobs
  end

  def create_resources
    case subject_type
    when "Organization"
      create_org_users_resource
      create_org_projects_resource
    when "OrganizationMembership"
      create_org_membership_resource
    when "Project"
      create_org_project_resource
    end
  end

  def create_org_users_resource
    resources.find_or_create_by!(resource_type: "users")
  end

  def create_org_projects_resource
    subject.projects.not_private.find_each do |project|
      resources.find_or_create_by!(resource_type: "project", resource_id: project.id)
      create_org_posts_resource(project)
      create_org_notes_resource(project)
      create_org_calls_resource(project)
    end
  end

  def create_org_posts_resource(project)
    project.kept_published_posts
      .eager_load(:attachments, kept_comments: :attachments)
      .find_each(batch_size: 50) do |post|
      resources.find_or_create_by!(resource_type: "post", resource_id: post.id)
      create_post_attachments_resource(post)
    end
  end

  def create_post_attachments_resource(post)
    attachments = post.attachments.to_a + post.kept_comments.flat_map(&:attachments)
    attachments.each do |attachment|
      resources.find_or_create_by!(resource_type: "attachment", resource_id: attachment.id)
    end
  end

  def create_org_notes_resource(project)
    project.kept_notes
      .eager_load(:attachments, kept_comments: :attachments)
      .find_each(batch_size: 50) do |note|
      resources.find_or_create_by!(resource_type: "note", resource_id: note.id)
      create_org_note_attachments_resource(note)
    end
  end

  def create_org_note_attachments_resource(note)
    attachments = note.attachments.to_a + note.kept_comments.flat_map(&:attachments)
    attachments.each do |attachment|
      resources.find_or_create_by!(resource_type: "attachment", resource_id: attachment.id)
    end
  end

  def create_org_calls_resource(project)
    project.calls
      .eager_load(:recordings)
      .find_each do |call|
      resources.find_or_create_by!(resource_type: "call", resource_id: call.id)
      create_org_call_recordings_resource(call)
    end
  end

  def create_org_call_recordings_resource(call)
    call.recordings.each do |recording|
      resources.find_or_create_by!(resource_type: "call_recording", resource_id: recording.id)
    end
  end

  def create_org_membership_resource
    resources.find_or_create_by!(resource_type: "member", resource_id: subject.id)

    projects = []

    subject.kept_published_posts
      .eager_load(:attachments, kept_comments: :attachments)
      .find_each(batch_size: 50) do |post|
      resources.find_or_create_by!(resource_type: "post", resource_id: post.id)
      create_post_attachments_resource(post)
      projects << post.project
    end

    projects.uniq.compact.each do |project|
      resources.find_or_create_by!(resource_type: "project", resource_id: project.id)
    end
  end

  def create_org_project_resource
    resources.find_or_create_by!(resource_type: "project", resource_id: subject.id)
    create_org_posts_resource(subject)
    create_org_notes_resource(subject)
    create_org_calls_resource(subject)
  end

  def queue_resource_jobs
    resources.find_each.with_index do |resource, index|
      DataExportResourceJob.perform_in(0.1.seconds * index, resource.id)
    end
  end

  def check_completed
    return if resources.pending.exists?

    Rails.logger.info("Data export #{public_id} completed, triggering task")
    run_task
  end

  def run_task
    # Replace ECS with local background job processing
    DataExportZipJob.perform_async(id)
  end

  def complete(zip_path)
    update!(zip_path: zip_path, completed_at: Time.current)

    OrganizationMailer.data_export_completed(self).deliver_later

    DataExportCleanupJob.perform_in(2.days, id)
  end

  def zip_url
    # Use local file serving instead of CloudFront
    Rails.application.routes.url_helpers.rails_blob_url(
      zip_path,
      host: Rails.application.config.action_mailer.default_url_options[:host]
    )
  end

  def cleanup!
    S3_BUCKET.object(zip_path).delete
    destroy!
  end

  # Local upload/export configuration
  def upload_name
    case subject_type
    when "Organization"
      "export-#{subject.slug}-#{public_id}"
    when "OrganizationMembership"
      "export-#{subject.user.username}-#{public_id}"
    else
      public_id
    end
  end
end
