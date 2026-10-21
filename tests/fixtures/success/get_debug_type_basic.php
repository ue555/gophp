<?php
class DebugSample {}
foreach ([null, true, 3, 2.5, 'text', [], new DebugSample, function () {}] as $value) {
    echo get_debug_type($value), '|';
}
