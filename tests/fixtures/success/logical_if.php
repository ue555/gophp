<?php
$a = 1;
$b = '1';
if ($a === 1 && $b === '1' || false) {
    echo 'outer|';
    if (!($a === $b) && !!$a) { echo 'inner|'; } else { echo 'bad'; }
} else { echo 'bad'; }
if (false && $missing || !true) { echo 'must not print'; } else { echo "else\n"; }
