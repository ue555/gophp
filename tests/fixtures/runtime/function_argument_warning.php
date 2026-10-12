<?php
function identity($v) { return $v; }
echo (identity($missing) === null) + 0, "\n";
