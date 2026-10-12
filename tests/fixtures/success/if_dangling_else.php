<?php
if (true)
    if (false) echo 'bad';
    else echo 'inner';
else echo 'bad';
if (false)
    if (true) echo 'bad';
    else echo 'bad';
else echo '-outer';
if (true) { if (false) echo 'bad'; } else echo 'bad';
echo "\n";
