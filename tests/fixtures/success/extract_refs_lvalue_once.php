<?php
class Bag { public $sets = [['left' => 1, 'right' => 2]]; }
$bag = new Bag;
$receiverCalls = 0;
$keyCalls = 0;
function receiver() { global $bag, $receiverCalls; $receiverCalls++; return $bag; }
function keyIndex() { global $keyCalls; $keyCalls++; return 0; }
echo extract(flags: EXTR_REFS, array: receiver()->sets[keyIndex()]), '|', $receiverCalls, '|', $keyCalls, '|';
$left = 9;
echo $bag->sets[0]['left'], '|';
$bag->sets[0]['right'] = 7;
echo $right;
