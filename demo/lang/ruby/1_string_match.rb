#!/usr/bin/env ruby
=begin
[case]
title=check string matches pattern
cid=0
pid=0

[group]
  1. exactly match >> hello
  2. regular expression match >> 1\d{10}
  3. format string match >> %s%d

[esac]
=end

print(">> hello\n");
print(">> 13905120512\n");
print(">> abc123\n");