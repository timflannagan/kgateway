/* eslint-disable */
/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var gogoproto_gogo_pb = require('../../../../gogoproto/gogo_pb.js');
var extproto_ext_pb = require('../../../../protoc-gen-ext/extproto/ext_pb.js');
goog.exportSymbol('proto.common.devportal.solo.io.ObjectMeta', null, global);
goog.exportSymbol('proto.common.devportal.solo.io.Time', null, global);

/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.common.devportal.solo.io.ObjectMeta = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.common.devportal.solo.io.ObjectMeta, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.common.devportal.solo.io.ObjectMeta.displayName = 'proto.common.devportal.solo.io.ObjectMeta';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.toObject = function(opt_includeInstance) {
  return proto.common.devportal.solo.io.ObjectMeta.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.common.devportal.solo.io.ObjectMeta} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.common.devportal.solo.io.ObjectMeta.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    namespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    uid: jspb.Message.getFieldWithDefault(msg, 5, ""),
    resourceversion: jspb.Message.getFieldWithDefault(msg, 6, ""),
    creationtimestamp: (f = msg.getCreationtimestamp()) && proto.common.devportal.solo.io.Time.toObject(includeInstance, f),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : [],
    annotationsMap: (f = msg.getAnnotationsMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.common.devportal.solo.io.ObjectMeta}
 */
proto.common.devportal.solo.io.ObjectMeta.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.common.devportal.solo.io.ObjectMeta;
  return proto.common.devportal.solo.io.ObjectMeta.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.common.devportal.solo.io.ObjectMeta} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.common.devportal.solo.io.ObjectMeta}
 */
proto.common.devportal.solo.io.ObjectMeta.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setNamespace(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setUid(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setResourceversion(value);
      break;
    case 8:
      var value = new proto.common.devportal.solo.io.Time;
      reader.readMessage(value,proto.common.devportal.solo.io.Time.deserializeBinaryFromReader);
      msg.setCreationtimestamp(value);
      break;
    case 11:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 12:
      var value = msg.getAnnotationsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.common.devportal.solo.io.ObjectMeta.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.common.devportal.solo.io.ObjectMeta} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.common.devportal.solo.io.ObjectMeta.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getUid();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getResourceversion();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getCreationtimestamp();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.common.devportal.solo.io.Time.serializeBinaryToWriter
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(11, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getAnnotationsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(12, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.common.devportal.solo.io.ObjectMeta.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string namespace = 3;
 * @return {string}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.getNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.common.devportal.solo.io.ObjectMeta.prototype.setNamespace = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string uid = 5;
 * @return {string}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.getUid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.common.devportal.solo.io.ObjectMeta.prototype.setUid = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string resourceVersion = 6;
 * @return {string}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.getResourceversion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.common.devportal.solo.io.ObjectMeta.prototype.setResourceversion = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional Time creationTimestamp = 8;
 * @return {?proto.common.devportal.solo.io.Time}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.getCreationtimestamp = function() {
  return /** @type{?proto.common.devportal.solo.io.Time} */ (
    jspb.Message.getWrapperField(this, proto.common.devportal.solo.io.Time, 8));
};


/** @param {?proto.common.devportal.solo.io.Time|undefined} value */
proto.common.devportal.solo.io.ObjectMeta.prototype.setCreationtimestamp = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.common.devportal.solo.io.ObjectMeta.prototype.clearCreationtimestamp = function() {
  this.setCreationtimestamp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.hasCreationtimestamp = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * map<string, string> labels = 11;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 11, opt_noLazyCreate,
      null));
};


proto.common.devportal.solo.io.ObjectMeta.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
};


/**
 * map<string, string> annotations = 12;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.common.devportal.solo.io.ObjectMeta.prototype.getAnnotationsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 12, opt_noLazyCreate,
      null));
};


proto.common.devportal.solo.io.ObjectMeta.prototype.clearAnnotationsMap = function() {
  this.getAnnotationsMap().clear();
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.common.devportal.solo.io.Time = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.common.devportal.solo.io.Time, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.common.devportal.solo.io.Time.displayName = 'proto.common.devportal.solo.io.Time';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.common.devportal.solo.io.Time.prototype.toObject = function(opt_includeInstance) {
  return proto.common.devportal.solo.io.Time.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.common.devportal.solo.io.Time} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.common.devportal.solo.io.Time.toObject = function(includeInstance, msg) {
  var f, obj = {
    seconds: jspb.Message.getFieldWithDefault(msg, 1, 0),
    nanos: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.common.devportal.solo.io.Time}
 */
proto.common.devportal.solo.io.Time.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.common.devportal.solo.io.Time;
  return proto.common.devportal.solo.io.Time.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.common.devportal.solo.io.Time} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.common.devportal.solo.io.Time}
 */
proto.common.devportal.solo.io.Time.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setSeconds(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setNanos(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.common.devportal.solo.io.Time.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.common.devportal.solo.io.Time.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.common.devportal.solo.io.Time} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.common.devportal.solo.io.Time.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSeconds();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getNanos();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
};


/**
 * optional int64 seconds = 1;
 * @return {number}
 */
proto.common.devportal.solo.io.Time.prototype.getSeconds = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.common.devportal.solo.io.Time.prototype.setSeconds = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional int32 nanos = 2;
 * @return {number}
 */
proto.common.devportal.solo.io.Time.prototype.getNanos = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.common.devportal.solo.io.Time.prototype.setNanos = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


goog.object.extend(exports, proto.common.devportal.solo.io);
