<?php
function out($v) { echo json_encode($v), "\n"; }
out(substr_replace('abcdef', 'X', 2, 2));
out(substr_replace('abcdef', 'X', -2));
out(substr_replace('abcdef', 'X', 2, -1));
out(substr_replace('abc', 'X', -99, 0));
out(substr_replace('abc', 'X', 99, 99));
out(substr_replace('abc', 'X', 1, -99));
out(substr_replace('日abc', 'Z', 3, 1));
out(substr_replace(['a'=>'abc', 4=>'def', 'ghi'], ['X','Y'], [1,2], [1]));
out(substr_replace('abc', ['X','Y'], 1, 1));
out(substr_replace(['abc','def'], 'X', 1, null));
foreach ([[[],1], [1,[]], [[1],[2,3]]] as [$start,$length]) {
 try { out(substr_replace('abc','X',$start,$length)); }
 catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
}
try { substr_replace('abc','X',new stdClass()); }
catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
