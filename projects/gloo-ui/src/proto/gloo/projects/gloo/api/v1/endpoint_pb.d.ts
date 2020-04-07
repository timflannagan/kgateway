/* eslint-disable */
// package: gloo.solo.io
// file: gloo/projects/gloo/api/v1/endpoint.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../protoc-gen-ext/extproto/ext_pb";
import * as solo_kit_api_v1_metadata_pb from "../../../../../solo-kit/api/v1/metadata_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../solo-kit/api/v1/ref_pb";
import * as solo_kit_api_v1_solo_kit_pb from "../../../../../solo-kit/api/v1/solo-kit_pb";

export class Endpoint extends jspb.Message {
  clearUpstreamsList(): void;
  getUpstreamsList(): Array<solo_kit_api_v1_ref_pb.ResourceRef>;
  setUpstreamsList(value: Array<solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addUpstreams(value?: solo_kit_api_v1_ref_pb.ResourceRef, index?: number): solo_kit_api_v1_ref_pb.ResourceRef;

  getAddress(): string;
  setAddress(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  getHostname(): string;
  setHostname(value: string): void;

  hasHealthCheck(): boolean;
  clearHealthCheck(): void;
  getHealthCheck(): HealthCheckConfig | undefined;
  setHealthCheck(value?: HealthCheckConfig): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: solo_kit_api_v1_metadata_pb.Metadata): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Endpoint.AsObject;
  static toObject(includeInstance: boolean, msg: Endpoint): Endpoint.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Endpoint, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Endpoint;
  static deserializeBinaryFromReader(message: Endpoint, reader: jspb.BinaryReader): Endpoint;
}

export namespace Endpoint {
  export type AsObject = {
    upstreamsList: Array<solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    address: string,
    port: number,
    hostname: string,
    healthCheck?: HealthCheckConfig.AsObject,
    metadata?: solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }
}

export class HealthCheckConfig extends jspb.Message {
  getHostname(): string;
  setHostname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthCheckConfig.AsObject;
  static toObject(includeInstance: boolean, msg: HealthCheckConfig): HealthCheckConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HealthCheckConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthCheckConfig;
  static deserializeBinaryFromReader(message: HealthCheckConfig, reader: jspb.BinaryReader): HealthCheckConfig;
}

export namespace HealthCheckConfig {
  export type AsObject = {
    hostname: string,
  }
}
