from pathlib import Path
import pg8000.dbapi
from dotenv import load_dotenv
import os


def execute_sql_file(cursor, filepath: Path):
    """
    Execute a SQL file.
    """
    if not filepath.exists():
        raise FileNotFoundError(f"SQL file not found: {filepath}")

    content = filepath.read_text(encoding="utf-8").strip()

    if not content:
        print(f"⚠️ Skipping empty file: {filepath.name}")
        return

    print(f"\n📄 Executing: {filepath.name}")
    # pg8000 maneja nativamente la ejecución de bloques de texto SQL
    cursor.execute(content)


def get_env():
    """
    Load environment variables from local .env
    """
    current_dir = Path(__file__).resolve().parent
    env_path = current_dir / ".env"

    if env_path.exists():
        load_dotenv(env_path, override=True)
        print(f"✅ Loaded .env from: {env_path}")
    else:
        print("⚠️ .env file not found, using fallback values")

    return {
        "host": os.getenv("DB_HOST_POSTGRES", "localhost"),
        # pg8000 requiere estrictamente que el puerto sea un entero
        "port": int(os.getenv("DB_PORT_POSTGRES", "5432")), 
        "user": os.getenv("DB_USER_POSTGRES", "postgres"),
        "password": os.getenv("DB_PASS_POSTGRES", "postgres"),
        # pg8000 usa 'database' en lugar de 'dbname'
        "database": os.getenv("DB_NAME_POSTGRES", "postgres"), 
    }


def show_menu():
    """
    Display menu options.
    """
    print("\n========== DATABASE MANAGER ==========")
    print("1️⃣  Execute schema.sql")
    print("2️⃣  Execute populate.sql")
    print("3️⃣  Execute BOTH")
    print("0️⃣  Exit")
    print("======================================")

    return input("Select an option: ").strip()


def main():
    current_dir = Path(__file__).resolve().parent

    db_config = get_env()

    option = show_menu()

    if option == "0":
        print("👋 Exiting...")
        return

    run_schema = option in ["1", "3"]
    run_populate = option in ["2", "3"]

    if not run_schema and not run_populate:
        print("❌ Invalid option")
        return

    conn = None
    try:
        print(
            f"\n🔌 Connecting to PostgreSQL "
            f"({db_config['host']}:{db_config['port']})..."
        )

        # Iniciar conexión usando DB-API de pg8000
        conn = pg8000.dbapi.connect(**db_config)
        cursor = conn.cursor()

        if run_schema:
            schema_path = current_dir.parent / "schema" / "schema.sql"
            execute_sql_file(cursor, schema_path)
            print("✅ schema.sql executed successfully")

        if run_populate:
            populate_path = current_dir.parent / "populate" / "populate.sql"
            execute_sql_file(cursor, populate_path)
            print("✅ populate.sql executed successfully")

        conn.commit()
        print("\n🎉 Database scripts executed successfully.")

    except FileNotFoundError as e:
        print(f"\n❌ {e}")

    except pg8000.dbapi.Error as e:
        print("\n❌ PostgreSQL error:")
        print(e)

    except Exception as e:
        print("\n❌ Unexpected error:")
        print(e)

    finally:
        if conn:
            conn.close()
            print("\n🔒 Database connection closed.")


if __name__ == "__main__":
    main()