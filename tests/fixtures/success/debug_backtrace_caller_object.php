<?php
function inspectCaller() {
    $frames = debug_backtrace(DEBUG_BACKTRACE_PROVIDE_OBJECT, 2);
    echo count($frames), '|', $frames[1]['function'], '|';
    echo $frames[1]['object'] === $GLOBALS['expectedCaller'] ? 'same' : 'different', '|';
    $without = debug_backtrace(0, 2);
    echo isset($without[1]['object']) ? 'present' : 'absent', '|';
}
class BacktraceCaller {
    public function inspect() { inspectCaller(); }
    public static function staticInspect() {
        $frames = debug_backtrace(DEBUG_BACKTRACE_PROVIDE_OBJECT, 2);
        echo isset($frames[0]['object']) ? 'present' : 'absent';
    }
}
$caller = new BacktraceCaller;
$GLOBALS['expectedCaller'] = $caller;
$caller->inspect();
BacktraceCaller::staticInspect();
