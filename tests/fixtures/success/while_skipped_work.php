<?php
while (false && $missing) { echo 1 / 0; }
$i = 0;
while ($i === 0 || 1 / 0) { $i++; break; echo 'bad'; }
echo $i, "\n";
