# frozen_string_literal: true

module Api
  module V1
    class UploadsController < BaseController
      before_action :authenticate_user!
      
      def create
        # Handle direct file uploads to replace S3 presigned posts
        file = params[:file]
        key = params[:key] || generate_upload_key(file)
        
        unless file.present?
          render json: { error: "No file provided" }, status: :bad_request
          return
        end
        
        # Validate file size and type if needed
        if file.size > 100.megabytes
          render json: { error: "File too large" }, status: :bad_request
          return
        end
        
        # Store the file using our local storage
        object = S3_BUCKET.object(key)
        object.put(body: file.tempfile)
        
        render json: {
          key: key,
          url: file_url(key),
          success: true
        }, status: :created
      end
      
      def show
        # Serve files from local storage
        key = params[:key]
        object = S3_BUCKET.object(key)
        
        unless object.exists?
          render json: { error: "File not found" }, status: :not_found
          return
        end
        
        send_file object.file_path, disposition: 'inline'
      end
      
      private
      
      def generate_upload_key(file)
        extension = File.extname(file.original_filename)
        timestamp = Time.current.to_i
        random = SecureRandom.hex(8)
        "uploads/#{timestamp}/#{random}#{extension}"
      end
      
      def file_url(key)
        url_for(controller: 'uploads', action: 'show', key: key)
      end
    end
  end
end