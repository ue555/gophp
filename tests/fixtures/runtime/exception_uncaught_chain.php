<?php
$first = new Exception('first');
echo "before\n";
throw new Error('second', 0, $first);
