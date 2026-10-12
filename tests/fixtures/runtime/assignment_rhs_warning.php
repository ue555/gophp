<?php
$a = 0;
echo true && ($a = $missing);
echo ($a === null) + 0, "\n";
