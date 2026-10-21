<?php
$path = tempnam(sys_get_temp_dir(), 'gophp-manifest-');
$replacement = tempnam(sys_get_temp_dir(), 'gophp-manifest-');
echo (int)($path !== $replacement), "\n";
echo (int)chmod($replacement, 0777 - umask()), "\n";
file_put_contents($replacement, '<?php return '.var_export(['providers' => ['Example']], true).';');
echo (int)rename($replacement, $path), "|", (int)file_exists($replacement), "\n";
$value = require $path;
echo $value['providers'][0], "\n";
echo (int)unlink($path), "|", (int)file_exists($path), "\n";
