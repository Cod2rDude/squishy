--!nolint
--!nocheck
--!optimize 2
--!native

--[[
    ******************************************************************************
    * @file     : ./tests/name.luau
    * @author   : squishy-compiler
    * @date     : February 1 2026
    * @lastEdit : February 1 2026 @ 21:52
    * @brief    : Squishy IDL Compiler generated code for name.
    * @version  : 1.0.0
    ******************************************************************************
    * @attention
    *
    * This software is licensed under terms that can be found in the LICENSE file 
    * in the root directory of this software component.
    * If no LICENSE file comes with this software, it is provided AS-IS.
    *
    ******************************************************************************
]]

--// Libs
local writer = require(script.Parent.Parent.libs.types.writer)
local reader = require(script.Parent.Parent.libs.types.reader)

--// Custom Type Definitions
type anotherStruct = {
    t0 : number;
}

--// Variables
local sharedBuffer = buffer.create(65536)

--// Functions
function write_anotherStruct(cursor : number, input : anotherStruct) : number
    cursor = writer.write_u8(sharedBuffer, cursor, input.t0)
    return cursor
end

function read_anotherStruct(buff : buffer, cursor : number) : (number, anotherStruct)
    local t0
    cursor, t0 = reader.read_u8(buff, cursor)
    return cursor, { t0 = t0; }
end

--// Lib Decleration
local scheme = {}

--// Lib Types
export type name = {
    t1 : number;
    t2 : number;
    t3 : number;
    t4 : number;
    t5 : number;
    t6 : number;
    t7 : number;
    t8 : number;
    t9 : number;
    t10 : Vector2;
    t11 : Vector3;
    t12 : Vector2int16;
    t13 : Vector3int16;
    t14 : Vector2;
    t15 : Vector3;
    t16 : CFrame;
    t17 : CFrame;
    t18 : CFrame;
    t19 : Color3;
    t20 : Color3;
    t22 : string;
    t23 : string;
    t24 : boolean;
    t25 : { [number] : boolean };
    t26 : { [number] : boolean };
    t27 : { [number] : number };
    t28 : { [number] : number };
    t29 : { [string] : number };
    t30 : { [string] : number };
    t31 : number;
    t32 : number;
    t33 : anotherStruct;
}

--// Lib Functions
function scheme.write(input : name) : buffer?
    local cursor = 2 -- Next is always at 2 at start because first 2 bytes are headers.
 
    cursor = writer.write_u8(sharedBuffer, cursor, input.t1)
    cursor = writer.write_i8(sharedBuffer, cursor, input.t2)
    cursor = writer.write_i16(sharedBuffer, cursor, input.t3)
    cursor = writer.write_u16(sharedBuffer, cursor, input.t4)
    cursor = writer.write_i16(sharedBuffer, cursor, input.t5)
    cursor = writer.write_u32(sharedBuffer, cursor, input.t6)
    cursor = writer.write_i32(sharedBuffer, cursor, input.t7)
    cursor = writer.write_f32(sharedBuffer, cursor, input.t8)
    cursor = writer.write_f64(sharedBuffer, cursor, input.t9)
    cursor = writer.write_vector2(sharedBuffer, cursor, input.t10)
    cursor = writer.write_vector3(sharedBuffer, cursor, input.t11)
    cursor = writer.write_vector2int16(sharedBuffer, cursor, input.t12)
    cursor = writer.write_vector3int16(sharedBuffer, cursor, input.t13)
    cursor = writer.write_vector2norm(sharedBuffer, cursor, input.t14)
    cursor = writer.write_vector3norm(sharedBuffer, cursor, input.t15)
    cursor = writer.write_cframe(sharedBuffer, cursor, input.t16)
    cursor = writer.write_cframe_e(sharedBuffer, cursor, input.t17)
    cursor = writer.write_cframe_q(sharedBuffer, cursor, input.t18)
    cursor = writer.write_color3(sharedBuffer, cursor, input.t19)
    cursor = writer.write_color3_hdr(sharedBuffer, cursor, input.t20)
    cursor = writer.write_string(sharedBuffer, cursor, input.t22)
    cursor = writer.write_string_l(sharedBuffer, cursor, input.t23)
    cursor = writer.write_bool(sharedBuffer, cursor, input.t24)
    cursor = writer.write_boolArray(sharedBuffer, cursor, input.t25, 16)
    cursor = writer.write_dynamicBoolArray(sharedBuffer, cursor, input.t26)
    cursor = writer.write_dynamicArray(sharedBuffer, cursor, input.t27, writer.write_u8)
    cursor = writer.write_array(sharedBuffer, cursor, input.t28, writer.write_u8, 16)
    cursor = writer.write_map(sharedBuffer, cursor, input.t29, writer.write_u8)
    cursor = writer.write_smap(sharedBuffer, cursor, input.t30, writer.write_u8)
    cursor = writer.write_f16(sharedBuffer, cursor, input.t31)
    cursor = writer.write_f24(sharedBuffer, cursor, input.t32)
    cursor = write_anotherStruct(cursor, input.t33)
 
    local packet = buffer.create(cursor)
    buffer.copy(packet, 2, sharedBuffer, 2, cursor-2)
    return packet
end

function scheme.read(buff : buffer) : name?
    local cursor = 2 -- Next is always at 2 at start because first 2 bytes are headers.
 
    local t1, t2, t3, t4, t5, t6, t7, t8, t9, t10, t11, t12, t13, t14, t15, t16, t17, t18, t19, t20, t22, t23, t24, t25, t26, t27, t28, t29, t30, t31, t32, t33
    cursor, t1 = reader.read_u8(buff, cursor)
    cursor, t2 = reader.read_i8(buff, cursor)
    cursor, t3 = reader.read_i16(buff, cursor)
    cursor, t4 = reader.read_u16(buff, cursor)
    cursor, t5 = reader.read_i16(buff, cursor)
    cursor, t6 = reader.read_u32(buff, cursor)
    cursor, t7 = reader.read_i32(buff, cursor)
    cursor, t8 = reader.read_f32(buff, cursor)
    cursor, t9 = reader.read_f64(buff, cursor)
    cursor, t10 = reader.read_vector2(buff, cursor)
    cursor, t11 = reader.read_vector3(buff, cursor)
    cursor, t12 = reader.read_vector2int16(buff, cursor)
    cursor, t13 = reader.read_vector3int16(buff, cursor)
    cursor, t14 = reader.read_vector2norm(buff, cursor)
    cursor, t15 = reader.read_vector3norm(buff, cursor)
    cursor, t16 = reader.read_cframe(buff, cursor)
    cursor, t17 = reader.read_cframe_e(buff, cursor)
    cursor, t18 = reader.read_cframe_q(buff, cursor)
    cursor, t19 = reader.read_color3(buff, cursor)
    cursor, t20 = reader.read_color3_hdr(buff, cursor)
    cursor, t22 = reader.read_string(buff, cursor)
    cursor, t23 = reader.read_string_l(buff, cursor)
    cursor, t24 = reader.read_bool(buff, cursor)
    cursor, t25 = reader.read_boolArray(sharedBuffer, cursor, 16)
    cursor, t26 = reader.read_dynamicBoolArray(sharedBuffer, cursor)
    cursor, t27 = reader.read_dynamicArray(buff, cursor, reader.read_u8)
    cursor, t28 = reader.read_array(buff, cursor, reader.read_u8, 16)
    cursor, t29 = reader.read_map(buff, cursor, reader.read_u8)
    cursor, t30 = reader.read_smap(buff, cursor, reader.read_u8)
    cursor, t31 = reader.read_f16(buff, cursor)
    cursor, t32 = reader.read_f24(buff, cursor)
    cursor, t33 = read_anotherStruct(buff, cursor)
 
    return { t1 = t1;
             t2 = t2;
             t3 = t3;
             t4 = t4;
             t5 = t5;
             t6 = t6;
             t7 = t7;
             t8 = t8;
             t9 = t9;
             t10 = t10;
             t11 = t11;
             t12 = t12;
             t13 = t13;
             t14 = t14;
             t15 = t15;
             t16 = t16;
             t17 = t17;
             t18 = t18;
             t19 = t19;
             t20 = t20;
             t22 = t22;
             t23 = t23;
             t24 = t24;
             t25 = t25;
             t26 = t26;
             t27 = t27;
             t28 = t28;
             t29 = t29;
             t30 = t30;
             t31 = t31;
             t32 = t32;
             t33 = t33;
             }
end

return scheme