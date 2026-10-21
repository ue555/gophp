<?php
class TraceBox { public $number = 1; }
function remember($numbers, $box) { return new Exception; }
$numbers = [4];
$box = new TraceBox;
$error = remember($numbers, $box);
$numbers[0] = 8;
$box->number = 2;
$trace = $error->getTrace();
echo $trace[0]['args'][0][0], '|';
echo ($trace[0]['args'][1] === $box ? 'same' : 'different'), '|';
echo $trace[0]['args'][1]->number, "\n";
