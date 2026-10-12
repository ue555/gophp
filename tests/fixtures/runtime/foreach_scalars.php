<?php
foreach ([null, false, true, 12, 1.5, 'abc'] as $input) {
    foreach ($input as $k => $v) echo 'bad';
    echo '.';
}
echo "done\n";
