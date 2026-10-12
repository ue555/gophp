<?php
foreach ([[], [0], ['x' => null]] as $v) echo (int)$v, ':', (float)$v, ':', (int)(bool)$v, ';';
echo "\n";
