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

var github_com_solo$io_solo$projects_projects_gloo_api_external_envoy_waf_waf_pb = require('../../../../../../../../../github.com/solo-io/solo-projects/projects/gloo/api/external/envoy/waf/waf_pb.js');
var gogoproto_gogo_pb = require('../../../../../../../../gogo/protobuf/gogoproto/gogo_pb.js');
goog.exportSymbol('proto.waf.plugins.gloo.solo.io.CoreRuleSet', null, global);
goog.exportSymbol('proto.waf.plugins.gloo.solo.io.RouteSettings', null, global);
goog.exportSymbol('proto.waf.plugins.gloo.solo.io.Settings', null, global);
goog.exportSymbol('proto.waf.plugins.gloo.solo.io.VhostSettings', null, global);

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
proto.waf.plugins.gloo.solo.io.Settings = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.waf.plugins.gloo.solo.io.Settings.repeatedFields_, null);
};
goog.inherits(proto.waf.plugins.gloo.solo.io.Settings, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.waf.plugins.gloo.solo.io.Settings.displayName = 'proto.waf.plugins.gloo.solo.io.Settings';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.waf.plugins.gloo.solo.io.Settings.repeatedFields_ = [3];



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
proto.waf.plugins.gloo.solo.io.Settings.prototype.toObject = function(opt_includeInstance) {
  return proto.waf.plugins.gloo.solo.io.Settings.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.waf.plugins.gloo.solo.io.Settings} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.Settings.toObject = function(includeInstance, msg) {
  var f, obj = {
    disabled: jspb.Message.getFieldWithDefault(msg, 1, false),
    coreRuleSet: (f = msg.getCoreRuleSet()) && proto.waf.plugins.gloo.solo.io.CoreRuleSet.toObject(includeInstance, f),
    ruleSetsList: jspb.Message.toObjectList(msg.getRuleSetsList(),
    github_com_solo$io_solo$projects_projects_gloo_api_external_envoy_waf_waf_pb.RuleSet.toObject, includeInstance)
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
 * @return {!proto.waf.plugins.gloo.solo.io.Settings}
 */
proto.waf.plugins.gloo.solo.io.Settings.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.waf.plugins.gloo.solo.io.Settings;
  return proto.waf.plugins.gloo.solo.io.Settings.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.waf.plugins.gloo.solo.io.Settings} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.waf.plugins.gloo.solo.io.Settings}
 */
proto.waf.plugins.gloo.solo.io.Settings.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 2:
      var value = new proto.waf.plugins.gloo.solo.io.CoreRuleSet;
      reader.readMessage(value,proto.waf.plugins.gloo.solo.io.CoreRuleSet.deserializeBinaryFromReader);
      msg.setCoreRuleSet(value);
      break;
    case 3:
      var value = new github_com_solo$io_solo$projects_projects_gloo_api_external_envoy_waf_waf_pb.RuleSet;
      reader.readMessage(value,github_com_solo$io_solo$projects_projects_gloo_api_external_envoy_waf_waf_pb.RuleSet.deserializeBinaryFromReader);
      msg.addRuleSets(value);
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
proto.waf.plugins.gloo.solo.io.Settings.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.waf.plugins.gloo.solo.io.Settings.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.waf.plugins.gloo.solo.io.Settings} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.Settings.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getCoreRuleSet();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.waf.plugins.gloo.solo.io.CoreRuleSet.serializeBinaryToWriter
    );
  }
  f = message.getRuleSetsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      github_com_solo$io_solo$projects_projects_gloo_api_external_envoy_waf_waf_pb.RuleSet.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool disabled = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.waf.plugins.gloo.solo.io.Settings.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.waf.plugins.gloo.solo.io.Settings.prototype.setDisabled = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional CoreRuleSet core_rule_set = 2;
 * @return {?proto.waf.plugins.gloo.solo.io.CoreRuleSet}
 */
proto.waf.plugins.gloo.solo.io.Settings.prototype.getCoreRuleSet = function() {
  return /** @type{?proto.waf.plugins.gloo.solo.io.CoreRuleSet} */ (
    jspb.Message.getWrapperField(this, proto.waf.plugins.gloo.solo.io.CoreRuleSet, 2));
};


/** @param {?proto.waf.plugins.gloo.solo.io.CoreRuleSet|undefined} value */
proto.waf.plugins.gloo.solo.io.Settings.prototype.setCoreRuleSet = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.waf.plugins.gloo.solo.io.Settings.prototype.clearCoreRuleSet = function() {
  this.setCoreRuleSet(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.waf.plugins.gloo.solo.io.Settings.prototype.hasCoreRuleSet = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * repeated envoy.config.filter.http.modsecurity.v2.RuleSet rule_sets = 3;
 * @return {!Array<!proto.envoy.config.filter.http.modsecurity.v2.RuleSet>}
 */
proto.waf.plugins.gloo.solo.io.Settings.prototype.getRuleSetsList = function() {
  return /** @type{!Array<!proto.envoy.config.filter.http.modsecurity.v2.RuleSet>} */ (
    jspb.Message.getRepeatedWrapperField(this, github_com_solo$io_solo$projects_projects_gloo_api_external_envoy_waf_waf_pb.RuleSet, 3));
};


/** @param {!Array<!proto.envoy.config.filter.http.modsecurity.v2.RuleSet>} value */
proto.waf.plugins.gloo.solo.io.Settings.prototype.setRuleSetsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.envoy.config.filter.http.modsecurity.v2.RuleSet=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.filter.http.modsecurity.v2.RuleSet}
 */
proto.waf.plugins.gloo.solo.io.Settings.prototype.addRuleSets = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.envoy.config.filter.http.modsecurity.v2.RuleSet, opt_index);
};


proto.waf.plugins.gloo.solo.io.Settings.prototype.clearRuleSetsList = function() {
  this.setRuleSetsList([]);
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
proto.waf.plugins.gloo.solo.io.CoreRuleSet = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.waf.plugins.gloo.solo.io.CoreRuleSet.oneofGroups_);
};
goog.inherits(proto.waf.plugins.gloo.solo.io.CoreRuleSet, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.waf.plugins.gloo.solo.io.CoreRuleSet.displayName = 'proto.waf.plugins.gloo.solo.io.CoreRuleSet';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.oneofGroups_ = [[2,3]];

/**
 * @enum {number}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.CustomsettingstypeCase = {
  CUSTOMSETTINGSTYPE_NOT_SET: 0,
  CUSTOM_SETTINGS_STRING: 2,
  CUSTOM_SETTINGS_FILE: 3
};

/**
 * @return {proto.waf.plugins.gloo.solo.io.CoreRuleSet.CustomsettingstypeCase}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.getCustomsettingstypeCase = function() {
  return /** @type {proto.waf.plugins.gloo.solo.io.CoreRuleSet.CustomsettingstypeCase} */(jspb.Message.computeOneofCase(this, proto.waf.plugins.gloo.solo.io.CoreRuleSet.oneofGroups_[0]));
};



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
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.toObject = function(opt_includeInstance) {
  return proto.waf.plugins.gloo.solo.io.CoreRuleSet.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.waf.plugins.gloo.solo.io.CoreRuleSet} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.toObject = function(includeInstance, msg) {
  var f, obj = {
    customSettingsString: jspb.Message.getFieldWithDefault(msg, 2, ""),
    customSettingsFile: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.waf.plugins.gloo.solo.io.CoreRuleSet}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.waf.plugins.gloo.solo.io.CoreRuleSet;
  return proto.waf.plugins.gloo.solo.io.CoreRuleSet.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.waf.plugins.gloo.solo.io.CoreRuleSet} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.waf.plugins.gloo.solo.io.CoreRuleSet}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCustomSettingsString(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCustomSettingsFile(value);
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
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.waf.plugins.gloo.solo.io.CoreRuleSet.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.waf.plugins.gloo.solo.io.CoreRuleSet} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {string} */ (jspb.Message.getField(message, 2));
  if (f != null) {
    writer.writeString(
      2,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string custom_settings_string = 2;
 * @return {string}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.getCustomSettingsString = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.setCustomSettingsString = function(value) {
  jspb.Message.setOneofField(this, 2, proto.waf.plugins.gloo.solo.io.CoreRuleSet.oneofGroups_[0], value);
};


proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.clearCustomSettingsString = function() {
  jspb.Message.setOneofField(this, 2, proto.waf.plugins.gloo.solo.io.CoreRuleSet.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.hasCustomSettingsString = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string custom_settings_file = 3;
 * @return {string}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.getCustomSettingsFile = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.setCustomSettingsFile = function(value) {
  jspb.Message.setOneofField(this, 3, proto.waf.plugins.gloo.solo.io.CoreRuleSet.oneofGroups_[0], value);
};


proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.clearCustomSettingsFile = function() {
  jspb.Message.setOneofField(this, 3, proto.waf.plugins.gloo.solo.io.CoreRuleSet.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.waf.plugins.gloo.solo.io.CoreRuleSet.prototype.hasCustomSettingsFile = function() {
  return jspb.Message.getField(this, 3) != null;
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
proto.waf.plugins.gloo.solo.io.VhostSettings = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.waf.plugins.gloo.solo.io.VhostSettings, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.waf.plugins.gloo.solo.io.VhostSettings.displayName = 'proto.waf.plugins.gloo.solo.io.VhostSettings';
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
proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.toObject = function(opt_includeInstance) {
  return proto.waf.plugins.gloo.solo.io.VhostSettings.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.waf.plugins.gloo.solo.io.VhostSettings} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.VhostSettings.toObject = function(includeInstance, msg) {
  var f, obj = {
    disabled: jspb.Message.getFieldWithDefault(msg, 1, false),
    settings: (f = msg.getSettings()) && proto.waf.plugins.gloo.solo.io.Settings.toObject(includeInstance, f)
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
 * @return {!proto.waf.plugins.gloo.solo.io.VhostSettings}
 */
proto.waf.plugins.gloo.solo.io.VhostSettings.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.waf.plugins.gloo.solo.io.VhostSettings;
  return proto.waf.plugins.gloo.solo.io.VhostSettings.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.waf.plugins.gloo.solo.io.VhostSettings} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.waf.plugins.gloo.solo.io.VhostSettings}
 */
proto.waf.plugins.gloo.solo.io.VhostSettings.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 2:
      var value = new proto.waf.plugins.gloo.solo.io.Settings;
      reader.readMessage(value,proto.waf.plugins.gloo.solo.io.Settings.deserializeBinaryFromReader);
      msg.setSettings(value);
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
proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.waf.plugins.gloo.solo.io.VhostSettings.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.waf.plugins.gloo.solo.io.VhostSettings} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.VhostSettings.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getSettings();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.waf.plugins.gloo.solo.io.Settings.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool disabled = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.setDisabled = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional Settings settings = 2;
 * @return {?proto.waf.plugins.gloo.solo.io.Settings}
 */
proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.getSettings = function() {
  return /** @type{?proto.waf.plugins.gloo.solo.io.Settings} */ (
    jspb.Message.getWrapperField(this, proto.waf.plugins.gloo.solo.io.Settings, 2));
};


/** @param {?proto.waf.plugins.gloo.solo.io.Settings|undefined} value */
proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.setSettings = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.clearSettings = function() {
  this.setSettings(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.waf.plugins.gloo.solo.io.VhostSettings.prototype.hasSettings = function() {
  return jspb.Message.getField(this, 2) != null;
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
proto.waf.plugins.gloo.solo.io.RouteSettings = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.waf.plugins.gloo.solo.io.RouteSettings, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.waf.plugins.gloo.solo.io.RouteSettings.displayName = 'proto.waf.plugins.gloo.solo.io.RouteSettings';
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
proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.toObject = function(opt_includeInstance) {
  return proto.waf.plugins.gloo.solo.io.RouteSettings.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.waf.plugins.gloo.solo.io.RouteSettings} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.RouteSettings.toObject = function(includeInstance, msg) {
  var f, obj = {
    disabled: jspb.Message.getFieldWithDefault(msg, 1, false),
    settings: (f = msg.getSettings()) && proto.waf.plugins.gloo.solo.io.Settings.toObject(includeInstance, f)
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
 * @return {!proto.waf.plugins.gloo.solo.io.RouteSettings}
 */
proto.waf.plugins.gloo.solo.io.RouteSettings.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.waf.plugins.gloo.solo.io.RouteSettings;
  return proto.waf.plugins.gloo.solo.io.RouteSettings.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.waf.plugins.gloo.solo.io.RouteSettings} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.waf.plugins.gloo.solo.io.RouteSettings}
 */
proto.waf.plugins.gloo.solo.io.RouteSettings.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 2:
      var value = new proto.waf.plugins.gloo.solo.io.Settings;
      reader.readMessage(value,proto.waf.plugins.gloo.solo.io.Settings.deserializeBinaryFromReader);
      msg.setSettings(value);
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
proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.waf.plugins.gloo.solo.io.RouteSettings.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.waf.plugins.gloo.solo.io.RouteSettings} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.waf.plugins.gloo.solo.io.RouteSettings.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getSettings();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.waf.plugins.gloo.solo.io.Settings.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool disabled = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.setDisabled = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional Settings settings = 2;
 * @return {?proto.waf.plugins.gloo.solo.io.Settings}
 */
proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.getSettings = function() {
  return /** @type{?proto.waf.plugins.gloo.solo.io.Settings} */ (
    jspb.Message.getWrapperField(this, proto.waf.plugins.gloo.solo.io.Settings, 2));
};


/** @param {?proto.waf.plugins.gloo.solo.io.Settings|undefined} value */
proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.setSettings = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.clearSettings = function() {
  this.setSettings(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.waf.plugins.gloo.solo.io.RouteSettings.prototype.hasSettings = function() {
  return jspb.Message.getField(this, 2) != null;
};


goog.object.extend(exports, proto.waf.plugins.gloo.solo.io);
