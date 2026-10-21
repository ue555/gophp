<?php
function show($v) { echo json_encode($v), "\n"; }
foreach ([["aaaaa", "aa"], ["日本日本", "日本"], ["banana", "an", 1], ["banana", "an", -5, -1], ["banana", "an", 2, 0], ["", "x"], [123123, 12]] as $args) { show(substr_count(...$args)); }
foreach ([["x", ""], ["x", "x", 2], ["abc", "a", 1, -3], ["abc", "a", 1, 3], [[], "a"], ["a", []], ["a", "a", []]] as $args) { try { substr_count(...$args); } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; } }
show(substr_count(needle: "a", haystack: "banana", length: null, offset: 2));
