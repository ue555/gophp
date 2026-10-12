<?php
for ($i = 0; $i < 0; ++$i) echo 'bad'; echo $i, ':';
for (; $i < 1; ) { echo $i; ++$i; } echo ':';
for ($i = 0; $i < 3; ++$i) echo $i;
for (;;) { echo ':end'; break; }
echo "\n";
