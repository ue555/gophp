<?php
echo ('01' === '1') + 0;
echo ('0e1' === '0e2') + 0;
echo ('10 cats' === '10 cats') + 0;
echo ('10 cats' !== 10) + 0;
echo ("A\0B" === "\x41\000B") + 0;
echo ("\0" === '') + 0;
echo ('a' === 'A') + 0;
echo ('日本語' === '日本語') + 0;
echo "\n";
