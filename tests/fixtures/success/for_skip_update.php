<?php
for ($i = 0; false; ++$missing) { echo 1 / 0; }
for (; true; 1 / 0) { break; }
for (; true; ++$missing) { break 1; }
echo 'done', "\n";
