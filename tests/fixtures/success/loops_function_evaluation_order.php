<?php
function mark($s, $v) { echo $s; return $v; }
for ($i = mark('I', 0); mark('C', $i < 2); mark('U', 0), ++$i) {
    echo 'B'; continue;
}
echo ':';
do { echo 'D'; } while (mark('Q', false));
echo ':';
while (mark('W', false)) { echo 'bad'; }
echo "\n";
