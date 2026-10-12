<?php
echo (!0) === true;
echo (!1) === false;
echo (!!'text') === true;
echo (1 && 'text') === true;
echo (0 && 'text') === false;
echo (0 || 'text') === true;
echo ('0' || null) === false;
echo ((1 && 1) === 1) === false;
echo ((0 || 0) === 0) === false;
echo ((0 || 'text') === 'text') === false;
echo "\n";
