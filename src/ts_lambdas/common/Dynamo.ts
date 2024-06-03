const { DynamoDBClient, PutItemCommand } = require('@aws-sdk/client-dynamodb');

const dynamoDBClient = new DynamoDBClient({ region: 'ap-southeast-1' });

const DynamoDB = {
    async write(hash, fileName, tableName) {
        const params = {
            TableName: tableName,
            Item: {
                hash: { S: hash },
                filename: { S: fileName }
            }
        };

        try {
            const command = new PutItemCommand(params);
            const newData = await dynamoDBClient.send(command);

            if (!newData) {
                throw new Error('There was an error writing the item');
            }

            return newData;
        } catch (error) {
            console.error(error);
            throw new Error(`Error writing item: ${error.message}`);
        }
    },
};

module.exports = DynamoDB;