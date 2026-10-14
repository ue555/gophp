<?php
echo "before\n";
function f(...$xs) { echo count($xs); }
if (false) {
    f(...[1], 2);
}
echo "after\n";
