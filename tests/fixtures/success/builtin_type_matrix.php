<?php
foreach ([null, false, 3, 1.5, '3', [], new stdClass] as $v) {
    echo gettype($v), ':', (int)is_int($v), (int)is_float($v), (int)is_string($v), (int)is_bool($v), (int)is_null($v), (int)is_array($v), (int)is_object($v), ';';
}
echo "\n", (int)IS_INT(3), ':', GETTYPE(3), "\n";
