/* eslint-disable */
// package: common.devportal.solo.io
// file: dev-portal/api/grpc/common/common.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class ObjectMeta extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getUid(): string;
  setUid(value: string): void;

  getResourceversion(): string;
  setResourceversion(value: string): void;

  hasCreationtimestamp(): boolean;
  clearCreationtimestamp(): void;
  getCreationtimestamp(): Time | undefined;
  setCreationtimestamp(value?: Time): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getAnnotationsMap(): jspb.Map<string, string>;
  clearAnnotationsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectMeta.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectMeta): ObjectMeta.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectMeta, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectMeta;
  static deserializeBinaryFromReader(message: ObjectMeta, reader: jspb.BinaryReader): ObjectMeta;
}

export namespace ObjectMeta {
  export type AsObject = {
    name: string,
    namespace: string,
    uid: string,
    resourceversion: string,
    creationtimestamp?: Time.AsObject,
    labelsMap: Array<[string, string]>,
    annotationsMap: Array<[string, string]>,
  }
}

export class Time extends jspb.Message {
  getSeconds(): number;
  setSeconds(value: number): void;

  getNanos(): number;
  setNanos(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Time.AsObject;
  static toObject(includeInstance: boolean, msg: Time): Time.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Time, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Time;
  static deserializeBinaryFromReader(message: Time, reader: jspb.BinaryReader): Time;
}

export namespace Time {
  export type AsObject = {
    seconds: number,
    nanos: number,
  }
}
