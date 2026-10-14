<?php

// Token-based candidate inventory. Review source and runtime paths before
// declaring a feature necessary or supported.
$root = __DIR__;
$lockPath = $root.'/composer.lock';
$vendorPath = $root.'/vendor';
if (!is_file($lockPath) || !is_dir($vendorPath)) {
    fwrite(STDERR, "composer.lock and vendor/ are required\n");
    exit(2);
}
$lock = json_decode(file_get_contents($lockPath), true, 512, JSON_THROW_ON_ERROR);
$packages = [];
foreach ($lock['packages'] ?? [] as $package) {
    $packages[$package['name']] = $package['version'];
}
$installedPath = $vendorPath.'/composer/installed.json';
if (!is_file($installedPath)) {
    fwrite(STDERR, "vendor/composer/installed.json is required\n");
    exit(2);
}
$installedData = json_decode(file_get_contents($installedPath), true, 512, JSON_THROW_ON_ERROR);
$installed = [];
foreach (($installedData['packages'] ?? $installedData) as $package) {
    $installed[$package['name']] = $package['version'];
}
foreach ($packages as $name => $version) {
    if (($installed[$name] ?? null) !== $version) {
        fwrite(STDERR, "lock/vendor version mismatch: $name ($version)\n");
        exit(2);
    }
}
$features = [
    T_ATTRIBUTE => 'attribute',
    T_FN => 'arrow_function',
    T_TRAIT => 'trait',
    T_INTERFACE => 'interface',
    T_ENUM => 'enum',
    T_MATCH => 'match',
    T_TRY => 'try',
    T_CATCH => 'catch',
    T_THROW => 'throw',
    T_INCLUDE => 'include',
    T_INCLUDE_ONCE => 'include_once',
    T_REQUIRE => 'require',
    T_REQUIRE_ONCE => 'require_once',
    T_YIELD => 'yield',
    T_YIELD_FROM => 'yield_from',
    T_READONLY => 'readonly',
    T_NULLSAFE_OBJECT_OPERATOR => 'nullsafe',
];
$inventory = [];
$iterator = new RecursiveIteratorIterator(new RecursiveDirectoryIterator($vendorPath));
foreach ($iterator as $file) {
    if (!$file->isFile() || $file->getExtension() !== 'php') {
        continue;
    }
    $relative = substr($file->getPathname(), strlen($vendorPath) + 1);
    $parts = explode('/', $relative);
    if (count($parts) < 3) {
        continue;
    }
    $name = $parts[0].'/'.$parts[1];
    if (!isset($packages[$name])) {
        continue;
    }
    foreach (token_get_all(file_get_contents($file->getPathname())) as $token) {
        if (!is_array($token) || !isset($features[$token[0]])) {
            continue;
        }
        $entry = &$inventory[$name][$features[$token[0]]];
        $entry['count'] = ($entry['count'] ?? 0) + 1;
        if (count($entry['examples'] ?? []) < 3) {
            $entry['examples'][] = 'vendor/'.$relative.':'.$token[2];
        }
        unset($entry);
    }
}
ksort($packages);
ksort($inventory);
echo json_encode(['packages' => $packages, 'syntax_candidates' => $inventory],
    JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_THROW_ON_ERROR), "\n";
