<?php
$a = ['x' => 1, 'y' => 2];
$b = ['y' => 9, 'z' => 3];
echo implode(',', array_keys(array_intersect_key($a, $b))), '|';
echo implode(',', array_flip(['a', 'b'])), '|';
$values = [3, 1, 2];
usort($values, static fn ($left, $right) => $left - $right);
echo implode(',', $values), '|';
echo preg_replace_callback('/(^|_)+(.)/', static fn ($m) => strtoupper($m[2]), 'foo_bar'), '|';
$object = new ArrayObject(['first' => 1]);
$object['second'] = 2;
echo $object['second'], ':', $object->count(), ':';
foreach ($object as $part) { echo $part; }
echo '|';
$text = 'abc';
$text[1] = 'Z';
echo $text, '|', sha1('abc'), '|', hash('xxh3', 'hello');
