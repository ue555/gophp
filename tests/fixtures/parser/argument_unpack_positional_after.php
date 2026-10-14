<?php
function f(...$xs) { echo count($xs); }
f(...[1], 2);
