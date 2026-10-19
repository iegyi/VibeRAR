use std::fs::File;
use std::io::Write;

fn main() {
    std::panic::set_hook(Box::new(|info| {
        let msg = match info.payload().downcast_ref::<&'static str>() {
            Some(s) => s.to_string(),
            None => match info.payload().downcast_ref::<String>() {
                Some(s) => s.clone(),
                None => "Box<Any>".to_string(),
            },
        };
        let loc = if let Some(location) = info.location() {
            format!("{}:{}:{}", location.file(), location.line(), location.column())
        } else {
            "unknown location".to_string()
        };
        if let Ok(mut f) = File::create("tauri_error.log") {
            let _ = writeln!(f, "Panic at {}: {}", loc, msg);
        }
    }));

    app_lib::run();
}
