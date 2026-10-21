<?php
$count = preg_match_all('/([a-z]+)=([0-9]+)/', 'a=1 b=22', $matches, PREG_SET_ORDER);
echo $count, '|', $matches[0][1], '|', $matches[1][2], '|';
preg_match_all('/b=([0-9]+)/', 'a=1 b=22', $offsets, PREG_OFFSET_CAPTURE);
echo $offsets[0][0][0], ':', $offsets[0][0][1], '|';
foreach (preg_split('/\r?\n/', "one\r\ntwo") as $part) {
    echo $part, '|';
}
