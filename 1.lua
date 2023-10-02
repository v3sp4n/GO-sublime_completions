
--SHIT_CODE
--SHIT_CODE
--SHIT_CODE

package.path = package.path .. ';../?.lua;../?.luac;../../?.lua;../../?.luac'
package.cpath = package.cpath .. ';../?.dll;../../?.dll'

req = require ('requests')

function removePattern(text)
    local p = '(){}|<>-.$^[]&*'
    for v in p:gmatch('.') do
        text = text:gsub('%'..v,'%%'..v)
    end
    return text
end

res = {}

URL = 'https://pkg.go.dev/time@go1.21.1'
TRIGGER = 'time.'

r = req.get(URL)
for l in r.text:gmatch('[^\n]+') do
    if l:find("%<pre%>func.+%<%/pre%>") then
        local func = l:match('%<pre%>(.-)<%/pre%>')
        if func ~= nil and not func:find('^func %(') then
            func = func:gsub('%<a href%=\".-\"%>',''):gsub('%<%/a%>','')
            ret = func:match('%)(.-)$')
            func = func:gsub(removePattern(ret)..'$','')
            a = func:match('%S+%s*%((.+)%)')
            if a ~= nil then
                func = func:gsub(removePattern(a),removePattern('${1:'..a..'}'))
            end
            print(func,"RET"..ret)
            table.insert(res,{
                trigger = TRIGGER .. func:match('^func (.-)%('),
                func = func,
                ret = ret,
                annotation = (a == nil and "" or a)
            })
        end
    end
end

j = ''

for k,v in pairs(res) do
    print(k,v.trigger,v.func)
    --                                              ${1:name}
    j = j .. '\n{\n\t"trigger": "' .. v.trigger .. '",\n\t"contents": "' .. TRIGGER .. v.func:gsub('^func%s*','') .. '",\n\t"details": "return ' .. v.ret .. '",\n\t"annotation": "('..v.annotation..')"\n},' 
end

local f = io.open("A.txt",'w')
f:write(j)
f:close()
