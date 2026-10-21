<?php
class TraceCalls {
    public function instance($item) { return new Exception('instance'); }
    public static function staticCall($item) { return new Error('static'); }
}
$object = new TraceCalls;
$first = $object->instance(2);
$second = TraceCalls::staticCall(3);
$closure = function ($item) { return new Exception('closure'); };
$third = $closure(4);
foreach ([$first, $second, $third] as $error) {
    $frame = $error->getTrace()[0];
    echo $frame['function'], '|';
    echo isset($frame['class']) ? $frame['class'] : '-', '|';
    echo isset($frame['type']) ? $frame['type'] : '-', '|';
    echo $frame['args'][0], "\n";
}
