<?php
function &sourceCell() {
    global $stored;
    return $stored;
}
function &relayCell() {
    return sourceCell();
}
$stored = 2;
$alias =& relayCell();
$alias = 6;
echo $stored, '|', sourceCell();
