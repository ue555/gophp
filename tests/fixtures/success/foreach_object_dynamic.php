<?php
$item = new stdClass;
$item->first = 1;
$item->second = [2, 3];
$count = 0;
foreach ($item as $key => $value) {
    echo $key, ':', gettype($value), ';';
    ++$count;
}
echo '|', $count, '|';
foreach (new stdClass as $value) echo 'bad';
echo 'done';
