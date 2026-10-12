<?php
$x = 'kept';
if (false) {
    $x = 'changed';
    echo 'bad';
    echo $missing;
    echo 1 / 0;
    echo 'text' + 1;
}
echo $x, "\n";
