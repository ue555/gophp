<?php
class HandlerBox {
    public static function staticHandle($level, $message) { echo "static:$level:$message|"; return true; }
    public function objectHandle($level, $message) { echo "object:$level:$message|"; return true; }
}
set_error_handler(['HandlerBox', 'staticHandle']);
trigger_error('one');
set_error_handler([new HandlerBox, 'objectHandle']);
trigger_error('two');
set_error_handler(function ($level, $message) { echo "closure:$level:$message|"; return true; });
trigger_error('three');
echo "done\n";
