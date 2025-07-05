# frozen_string_literal: true

# Local file storage replacement for AWS S3
class LocalStorage
  attr_reader :base_path

  def initialize(base_path = nil)
    @base_path = base_path || Rails.root.join("storage", "uploads")
    FileUtils.mkdir_p(@base_path)
  end

  def bucket(name)
    LocalBucket.new(File.join(@base_path, name))
  end
end

class LocalBucket
  attr_reader :path

  def initialize(path)
    @path = path
    FileUtils.mkdir_p(@path)
  end

  def object(key)
    LocalObject.new(File.join(@path, key))
  end

  def presigned_post(options = {})
    LocalPresignedPost.new(self, options)
  end
end

class LocalObject
  attr_reader :file_path

  def initialize(file_path)
    @file_path = file_path
  end

  def key
    File.basename(@file_path)
  end

  def put(body:, **options)
    FileUtils.mkdir_p(File.dirname(@file_path))
    if body.respond_to?(:read)
      File.open(@file_path, 'wb') { |f| IO.copy_stream(body, f) }
    else
      File.write(@file_path, body)
    end
    self
  end

  def copy_to(target_object)
    FileUtils.mkdir_p(File.dirname(target_object.file_path))
    FileUtils.cp(@file_path, target_object.file_path)
    target_object
  end

  def delete
    File.delete(@file_path) if File.exist?(@file_path)
    self
  end

  def exists?
    File.exist?(@file_path)
  end

  def read
    File.read(@file_path)
  end
end

class LocalPresignedPost
  attr_reader :bucket, :options, :url, :fields

  def initialize(bucket, options = {})
    @bucket = bucket
    @options = options
    @url = "/api/v1/uploads" # We'll need to create this endpoint
    @fields = {
      "key" => options[:key],
      "Content-Type" => options[:content_type],
      "success_action_status" => options[:success_action_status] || "201",
      "policy" => "local-policy-placeholder",
      "x-amz-algorithm" => "AWS4-HMAC-SHA256",
      "x-amz-credential" => "local-credential",
      "x-amz-date" => Time.current.utc.strftime("%Y%m%dT%H%M%SZ"),
      "x-amz-signature" => "local-signature"
    }
  end
end

# Exception classes to match AWS S3 API
module LocalStorageErrors
  class NoSuchKey < StandardError; end
end

# Initialize the local storage
LOCAL_STORAGE = LocalStorage.new
S3_BUCKET = LOCAL_STORAGE.bucket(Rails.env.production? ? "campsite-media" : "campsite-media-dev")