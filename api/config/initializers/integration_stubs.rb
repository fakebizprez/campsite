# frozen_string_literal: true

# Stub implementations for external integrations that replace cloud services
# These can be extended or replaced with open-source alternatives

module IntegrationStubs
  # Slack integration stub
  class SlackStub
    def self.enabled?
      false
    end

    def self.send_notification(message, channel = nil)
      Rails.logger.info("SLACK STUB: Would send '#{message}' to channel '#{channel}'")
      { success: true, message: "Slack integration disabled - using stub" }
    end

    def self.create_channel(name)
      Rails.logger.info("SLACK STUB: Would create channel '#{name}'")
      { id: "stub_#{SecureRandom.hex(4)}", name: name }
    end
  end

  # Linear integration stub  
  class LinearStub
    def self.enabled?
      false
    end

    def self.create_issue(title, description)
      Rails.logger.info("LINEAR STUB: Would create issue '#{title}'")
      { id: "stub_#{SecureRandom.hex(4)}", title: title, url: "https://linear.stub/issue/123" }
    end

    def self.sync_comments(issue_id, comments)
      Rails.logger.info("LINEAR STUB: Would sync #{comments.length} comments to issue #{issue_id}")
      { synced: comments.length }
    end
  end

  # Figma integration stub
  class FigmaStub
    def self.enabled?
      false
    end

    def self.get_file_info(file_key)
      Rails.logger.info("FIGMA STUB: Would get file info for '#{file_key}'")
      { 
        name: "Stub Design File",
        thumbnail_url: "/placeholder-thumbnail.png",
        last_modified: Time.current.iso8601
      }
    end

    def self.render_frame(file_key, node_id)
      Rails.logger.info("FIGMA STUB: Would render frame '#{node_id}' from file '#{file_key}'")
      { image_url: "/placeholder-frame.png" }
    end
  end

  # Cal.com integration stub
  class CalDotComStub
    def self.enabled?
      false
    end

    def self.create_booking(event_type, datetime)
      Rails.logger.info("CAL.COM STUB: Would create booking for #{event_type} at #{datetime}")
      { 
        id: "stub_#{SecureRandom.hex(4)}", 
        url: "https://cal.stub/booking/123",
        status: "confirmed"
      }
    end
  end

  # Email service stub (replaces Postmark)
  class EmailStub
    def self.enabled?
      false
    end

    def self.send_email(to:, subject:, body:, template: nil)
      Rails.logger.info("EMAIL STUB: Would send email to '#{to}' with subject '#{subject}'")
      if Rails.env.development?
        puts "\n" + "="*50
        puts "EMAIL STUB"
        puts "To: #{to}"
        puts "Subject: #{subject}"
        puts "Template: #{template}" if template
        puts "-" * 20
        puts body
        puts "="*50 + "\n"
      end
      { message_id: "stub_#{SecureRandom.hex(8)}", status: "sent" }
    end
  end

  # Transcription service stub (replaces AWS Transcribe)
  class TranscriptionStub
    def self.enabled?
      false
    end

    def self.transcribe_audio(audio_url)
      Rails.logger.info("TRANSCRIPTION STUB: Would transcribe audio from '#{audio_url}'")
      # Return placeholder transcription
      {
        transcript: "[Transcription service disabled - placeholder text]",
        confidence: 0.0,
        status: "completed"
      }
    end

    def self.generate_vtt(transcript_data)
      # Return basic VTT format for captions
      <<~VTT
        WEBVTT

        00:00:00.000 --> 00:05:00.000
        [Transcription service disabled - placeholder text]
      VTT
    end
  end

  # Search service stub (can be replaced with local search like Elasticsearch)
  class SearchStub
    def self.enabled?
      false
    end

    def self.index_document(type, id, content)
      Rails.logger.info("SEARCH STUB: Would index #{type}:#{id}")
      true
    end

    def self.search(query, filters = {})
      Rails.logger.info("SEARCH STUB: Would search for '#{query}'")
      { results: [], total: 0, message: "Search service disabled" }
    end

    def self.delete_document(type, id)
      Rails.logger.info("SEARCH STUB: Would delete #{type}:#{id}")
      true
    end
  end
end

# Make stubs available globally
SlackService = IntegrationStubs::SlackStub
LinearService = IntegrationStubs::LinearStub  
FigmaService = IntegrationStubs::FigmaStub
CalDotComService = IntegrationStubs::CalDotComStub
EmailService = IntegrationStubs::EmailStub
TranscriptionService = IntegrationStubs::TranscriptionStub
SearchService = IntegrationStubs::SearchStub