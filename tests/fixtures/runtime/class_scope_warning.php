<?php
class LocalWarning {
    function read() { return $global; }
}
$global = 7; $o = new LocalWarning;
echo ($o->read() === null) + 0, ':', $global, "\n";
