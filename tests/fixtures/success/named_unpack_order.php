<?php
function show($a, $b, $c, ...$rest) {
    echo $a, $b, $c, '|';
    foreach ($rest as $key => $value) echo $key, '=', $value, ';';
    echo "\n";
}
$i = 0;
function nextValue() { global $i; $i++; return $i; }
show(...[nextValue(), 'c' => nextValue()], b: nextValue(), extra: nextValue());
echo $i, "\n";
