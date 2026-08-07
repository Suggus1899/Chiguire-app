mod fiscal;

use tauri::Manager;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(fiscal::plugin::init())
        .setup(|app| {
            #[cfg(debug_assertions)]
            {
                if let Some(window) = app.get_webview_window("main") {
                    window.open_devtools();
                }
            }
            log::info!("Chiguire ERP desktop starting");
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running Chiguire ERP");
}
