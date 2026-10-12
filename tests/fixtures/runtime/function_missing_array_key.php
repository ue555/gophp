<?php
function read($a) { return $a['missing']; }
echo (read([1]) === null) + 0, "\n";
