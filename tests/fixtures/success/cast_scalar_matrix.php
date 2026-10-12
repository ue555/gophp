<?php
foreach ([null, false, true, 0, 3, -2, 2.75, -0.0, '', '0', '12.5tail', 'bad'] as $v) {
    echo (int)(bool)$v, ':', (int)$v, ':', (float)$v, ':', (string)$v, ';';
}
echo "\n";
