<?php
class TraceOrigin {
    public function show() {
        $frame = debug_backtrace(DEBUG_BACKTRACE_PROVIDE_OBJECT, 1)[0];
        echo $frame['class'], '|', $frame['object'] === $this ? 'same' : 'different', '|';
        echo implode(',', array_keys($frame));
    }
}
class TraceChild extends TraceOrigin {}
(new TraceChild)->show();
