<?php
$process = proc_open('printf hello', [0 => ['pipe', 'r'], 1 => ['pipe', 'w'], 2 => ['pipe', 'w']], $pipes);
fclose($pipes[0]);
$read = [$pipes[1]];
$write = $except = [];
echo stream_select($read, $write, $except, 1) > 0 ? 'ready|' : 'bad|';
echo stream_get_contents($pipes[1]), '|';
fclose($pipes[1]);
fclose($pipes[2]);
echo proc_close($process);
