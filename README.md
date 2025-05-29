# ActionAI

ActionAI is a **prototype program** that allows you to leverage context-aware actions powered by AI. It's a functional experiment and a work in progress. It is being developed in GNOME desktop environment and currently supports Wayland.

## Project Status

**This is a prototype program.** It is functional and works, but it's primarily an experiment and a work in progress.

## Features

*   **AI-Powered Actions:** Execute predefined and custom actions using advanced AI models.
*   **Context-Aware Integration:** Leverage desktop context like selected text, clipboard content, and screenshots for intelligent AI interactions.
*   **Multi-Modal Inputs:** Combine multiple input types such as voice, screen sections, and text for rich, dynamic prompts.
*   **Flexible Outputs:** Deliver AI responses through various channels, including dialogs, clipboard, voice playback, and console output.
*   **Customizable Actions:** Define declarative triggers for customizable actions. See [`configs/shortcuts.yaml`](configs/shortcuts.yaml) as an example.

## Practical Use Cases

Here are some practical examples to inspire your usage:

1.  **Correct Selected Text:**
    *   **Objective:** Select text on your screen and let the AI correct its grammar and style. The corrected version is then copied to the clipboard.
    *   **Command:**
        ```bash
        actionai run -i selected-text -o clipboard -n "Correct this text. Reply only with the corrected text."
        ```

2.  **Explain a Screen Section:**
    *   **Objective:** Capture a section of your screen and explain it. The response is displayed in a dialog window.
    *   **Command:**
        ```bash
        actionai run -i screen-section -o window -n "Explain this image"
        ```
3.  **Ask a Question Using Voice:**
    *   **Objective:** Use your voice to ask a question and receive a spoken answer from the AI.
    *   **Command:**
        ```bash
        actionai run -i voice -o voice
        ```

4.  **Combine Screen Section and Voice Input:**
    *   **Objective:** Capture a part of your screen and then ask a question about it using your voice. The AI's answer is displayed in a dialog window.
    *   **Command:**
        ```bash
        actionai run -i screen-section,voice -o window
        ```

## Context-Aware Capabilities

ActionAI acts on your current desktop context (GNOME environment) to provide relevant AI assistance, integrating with your workflow.

### Inputs & Outputs

**Input:**
*   `clipboard`: Uses content from the system clipboard.
*   `screen`: Captures the entire screen.
*   `screen-section`: Captures a selected section of the screen.
*   `selected-text`: Captures the currently selected text.
*   `voice`: Uses voice as input.
*   `window`: Prompts for input in a dialog.

**Output:**
*   `clipboard`: Places the output onto the system clipboard.
*   `stdout`: Prints the output to the standard console.
*   `voice`: Communicates output with voice.
*   `window`: Displays the output in a dialog.

## Usage

The basic command structure is:

```bash
actionai run -m <model_id> -i <input_type(s)> -o <output_type(s)> [-n "<system-instructions>"]
```

*   `-m <model_id>`: Specifies the AI model to use (e.g., `gpt-4.1-mini`).
*   `-i <input_type(s)>`: Defines one or more input sources. Multiple inputs can be comma-separated (e.g., `-i screen-section,voice`). The order matters, as the data from these inputs will be concatenated to form the final prompt for the AI.
*   `-o <output_type(s)>`: Defines one or more output destinations. Multiple outputs can be comma-separated (e.g., `-o window,clipboard`). The AI's response will be sent to all specified outputs.
*   `-n "<system-instructions>"`: (Optional) A specific instruction or question for the AI. If not provided, the AI will act based on the combined context from the input sources.

### Multi-modal input pipeline

You can chain multiple input types. For example, `-i screen-section,voice` will first capture a section of the screen, then record voice input. Both pieces of information are then combined and sent to the AI. This allows for rich, multi-modal interactions.

### Output

Similarly, you can direct the AI's response to the desired output. For instance, `-o window` will display the response in a dialog window, while `-o clipboard` will send to the clipboard. 

## Prerequisites

*   For GNOME-based systems, the following command-line utilities are required:
    *   `aplay`: Used for audio playback. (Typically part of `alsa-utils` package)
    *   `zenity`: Used to display native GTK dialogs for information and input.
    *   `ffmpeg`: Used for voice recording and audio format conversion.
    *   `notify-send`: Used for sending desktop notifications.
    *   `gnome-screenshot`: Used for taking screenshots on GNOME desktops. 
    *   `wl-clipboard` (provides `wl-copy` and `wl-paste`): Used for clipboard operations (copy/paste) and accessing the primary selection on Wayland.

## Integration with GNOME Custom Shortcuts

ActionAI is designed to be used as a backend, and no frontend is currently available. However, you can leverage GNOME's custom shortcuts to quickly trigger preconfigured AI actions. This allows you to integrate ActionAI seamlessly into your workflow without needing a dedicated graphical interface. Check the [`configs/shortcuts.yaml`](configs/shortcuts.yaml) file for declarative action definitions.

By combining ActionAI with GNOME custom shortcuts, you can create a personalized and efficient AI-powered workflow.

## License

This project is licensed under the terms of the `LICENSE` file.
