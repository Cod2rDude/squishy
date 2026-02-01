--!nolint
--!nocheck
--!optimize 2
--!native

--[[
    ******************************************************************************
    * @file     : squishy/squishy-compiler/internal/app/assets/template.luau      
    * @author   : Cod2rDude
    * @date     : January 24 2026
    * @lastEdit : January 24 2025 @ 22:24
    * @brief    : Squishy-Compiler output template.
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

--// Variables
local sharedBuffer = buffer.create(65536)

--// Functions

--// Lib Decleration
local scheme = {}

--// Lib Types

--// Lib Functions
function scheme.write(input : any) : buffer?
    local cursor = 2 -- Next is always at 2 at start because first 2 bytes are headers.

    local packet = buffer.create(cursor)
    buffer.copy(packet, 2, sharedBuffer, 2, cursor-2)
    return packet
end

function scheme.read(buff : buffer) : any?
    local cursor = 2 -- Next is always at 2 at start because first 2 bytes are headers.

    return 
end

return scheme