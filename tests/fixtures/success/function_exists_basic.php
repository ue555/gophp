<?php
function customFunction() { return 1; }
echo function_exists('CUSTOMFUNCTION') ? 'yes|' : 'no|';
echo function_exists('strlen') ? 'builtin|' : 'missing|';
echo function_exists('not_defined_here') ? 'wrong' : 'clear';
