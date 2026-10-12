<?php
echo (1 + 2 * 3 === 7) + 0;
echo ('sum=' . 2 + 3 === 'sum=5') + 0;
echo ('x' . 1 !== 'x1') + 0;
echo (-2 * 3 === -6) + 0;
echo ((1 === 1) === true) + 0;
echo (1 === (1 !== 2)) + 0;
echo (true !== (1 === 2)) + 0;
echo "\n";
