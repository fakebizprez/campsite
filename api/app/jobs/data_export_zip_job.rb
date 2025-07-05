# frozen_string_literal: true

# Local replacement for AWS ECS data export task
class DataExportZipJob < BaseJob
  sidekiq_options queue: "background"

  def perform(data_export_id)
    data_export = DataExport.find(data_export_id)
    
    Rails.logger.info("Starting local data export zip creation for #{data_export.public_id}")
    
    begin
      zip_path = create_zip_archive(data_export)
      data_export.complete(zip_path)
      Rails.logger.info("Data export #{data_export.public_id} completed successfully")
    rescue => e
      Rails.logger.error("Data export #{data_export.public_id} failed: #{e.message}")
      Sentry.capture_exception(e, tags: { data_export_id: data_export.id })
      raise
    end
  end

  private

  def create_zip_archive(data_export)
    export_dir = File.join(LOCAL_STORAGE.base_path, S3_BUCKET.path, "exports", data_export.public_id)
    zip_filename = "#{data_export.upload_name}.zip"
    zip_path = "exports/#{data_export.public_id}/#{zip_filename}"
    full_zip_path = File.join(LOCAL_STORAGE.base_path, S3_BUCKET.path, zip_path)
    
    # Create zip directory if it doesn't exist
    FileUtils.mkdir_p(File.dirname(full_zip_path))
    
    # Create the zip file
    require 'zip'
    
    Zip::File.open(full_zip_path, Zip::File::CREATE) do |zipfile|
      if Dir.exist?(export_dir)
        Dir.glob(File.join(export_dir, "**", "*")).each do |file|
          next if File.directory?(file)
          
          # Calculate the path relative to the export directory
          relative_path = file.sub("#{export_dir}/", "")
          zipfile.add(relative_path, file)
        end
      end
    end
    
    zip_path
  end
end